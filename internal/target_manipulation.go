package deglob

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

type Target struct {
	start          int
	end            int
	name           string
	content        []string
	glob_attr      string
	glob_attr_line int
	glob_patterns  []string
	globbed_files  []string
}

var glob_pattern = regexp.MustCompile(`\s(?P<attr>.*)\s=\sglob\(\[(?P<files>.*)\]\)`)
var target_name_pattern = regexp.MustCompile(`name\s=\s\"(?P<name>.*)\"`)
var target_start_pattern = regexp.MustCompile(`^cc_library\(.*`) // TODO: Support cc_binary
var target_end_pattern = regexp.MustCompile(`^\)\n$`)

func ExtractTargetsFromFileContents(file_contents []string, filtered_file string) []Target {
	var targets_with_glob []Target

	var current_target_glob_patterns []string
	var current_target_content []string
	current_target_name := ""
	current_target_glob_attr := ""
	current_target_glob_attr_line_number := -1
	currently_in_target := false
	current_target_start_line_number := -1

	current_line_number := 0
	for _, line := range file_contents {
		if target_start_pattern.MatchString(line) && !currently_in_target {
			currently_in_target = true
			current_target_start_line_number = current_line_number
		}

		if currently_in_target {
			current_target_content = append(current_target_content, line)
		}

		if currently_in_target && target_name_pattern.MatchString(line) {
			matches := target_name_pattern.FindStringSubmatch(line)
			name_index := target_name_pattern.SubexpIndex("name")
			current_target_name = matches[name_index]
		}

		if currently_in_target && target_end_pattern.MatchString(line) {
			// Only track target if the target had a glob
			if len(current_target_glob_patterns) > 0 {
				target := Target{start: current_target_start_line_number, end: current_line_number, name: current_target_name, glob_patterns: append([]string(nil), current_target_glob_patterns...), content: append([]string(nil), current_target_content...), glob_attr: current_target_glob_attr, glob_attr_line: current_target_glob_attr_line_number} // deep copy the slices
				target.globbed_files = find_files_from_glob_in_targets(target, filtered_file)

				targets_with_glob = append(targets_with_glob, target)

				fmt.Println("Found glob in target: " + current_target_name + " with attr " + current_target_glob_attr + " on line " + fmt.Sprint(current_target_glob_attr_line_number))
			}

			currently_in_target = false
			clear(current_target_glob_patterns)
			clear(current_target_content)
		}

		// 4) Find any globs in the files and the patterns they capture
		if currently_in_target && glob_pattern.MatchString(line) {
			matches := glob_pattern.FindStringSubmatch(line)
			files_index := glob_pattern.SubexpIndex("files")
			attr_index := glob_pattern.SubexpIndex("attr")

			// TODO: Tackle multiple strings inside a single glob
			glob_pattern := strings.Trim(matches[files_index], "\"")
			current_target_glob_attr = strings.TrimSpace(matches[attr_index])
			current_target_glob_patterns = append(current_target_glob_patterns, glob_pattern)
			current_target_glob_attr_line_number = current_line_number
		}

		current_line_number += 1
	}

	return targets_with_glob
}

func CreateNewFileContentsIncludingNewTargets(existing_file_contents []string, targets_with_glob []Target) []string {
	var new_file_contents []string

	currently_in_target := false
	var new_target_content []string
	var current_target Target

	for index, line := range existing_file_contents {
		line_to_add := line
		if currently_in_target {
			if index == current_target.glob_attr_line {
				line_to_add = strings.Replace(line, current_target.glob_attr, "deps", -1)
				line_to_add = strings.Replace(line_to_add, "glob([\""+current_target.glob_patterns[0]+"\"])", "["+create_list_of_new_target_names_from_target(current_target)+"]", -1)
			}
			if index == current_target.end {
				currently_in_target = false
				new_target_content = create_new_targets_from_globbed_files(current_target)
			}
		} else if len(new_target_content) > 0 {
			new_file_contents = append(new_file_contents, new_target_content...)
			clear(new_target_content)
			new_target_content = nil
		} else if !currently_in_target && is_current_line_number_within_any_target(index, targets_with_glob) {
			currently_in_target = true
			current_target = return_target_in_current_line_number(index, targets_with_glob)
		}

		new_file_contents = append(new_file_contents, line_to_add)
	}

	return new_file_contents
}

func find_files_from_glob_in_targets(target Target, filtered_file string) []string {
	package_path := strings.Replace(filtered_file, "BUILD", "", -1)

	fmt.Println("Package_path: ", package_path)
	fmt.Println("Target content: ", target.content)
	fmt.Println("Glob attr: ", target.glob_attr)

	var globbed_files []string

	for _, glob_pattern := range target.glob_patterns {
		glob_cmd := strings.Replace(filtered_file, "BUILD", glob_pattern, -1)
		fmt.Println("Glob_cmd: ", glob_cmd)
		files, _ := filepath.Glob(glob_cmd)

		for _, file := range files {
			globbed_full_file_path := package_path + file
			globbed_package_relative_file_path := strings.Replace(globbed_full_file_path, package_path, "", -1)
			fmt.Println("Relative_file_path: ", globbed_package_relative_file_path)
			globbed_files = append(globbed_files, globbed_package_relative_file_path)
		}
	}

	return globbed_files
}

func create_new_targets_from_globbed_files(target Target) []string {
	var new_target_content []string
	for _, target_globbed_file := range target.globbed_files {
		new_target_content = append(new_target_content, "\n")
		for _, target_content_line := range target.content {
			if glob_pattern.MatchString(target_content_line) {
				// If this is the glob line of the target, replace it with the explicit source file
				// TODO: Work with multiple glob patterns
				new_src_line := strings.Replace(target_content_line, "glob([\""+target.glob_patterns[0]+"\"])", "[\""+target_globbed_file+"\"]", -1)
				new_target_content = append(new_target_content, new_src_line)
			} else if target_name_pattern.MatchString(target_content_line) {
				// If this is the name line of the target, replace it with a new name
				new_name_suffix := strings.Split(target_globbed_file, ".")[0]
				new_name_line := strings.Replace(target_content_line, target.name, target.name+"_"+new_name_suffix, -1)
				new_target_content = append(new_target_content, new_name_line)
			} else {
				new_target_content = append(new_target_content, target_content_line)
			}
		}
	}
	return new_target_content
}

func create_list_of_new_target_names_from_target(target Target) string {
	var new_target_names []string
	for _, target_globbed_file := range target.globbed_files {
		for _, target_content_line := range target.content {
			if target_name_pattern.MatchString(target_content_line) {
				// If this is the name line of the target, replace it with a new name
				new_name_suffix := strings.Split(target_globbed_file, ".")[0]
				new_target_names = append(new_target_names, "\":"+target.name+"_"+new_name_suffix+"\"")
			}
		}
	}
	return strings.Join(new_target_names, ", ")
}
