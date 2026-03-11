cc_library(
    name = "files_with_glob",
    deps = [":files_with_glob_file_1_cpp", ":files_with_glob_file_2_cpp"],
)

cc_library(
    name = "files_with_glob_file_1_cpp",
    hdrs = ["file_1.cpp"],
)

cc_library(
    name = "files_with_glob_file_2_cpp",
    hdrs = ["file_2.cpp"],
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":files_with_glob"],
)
