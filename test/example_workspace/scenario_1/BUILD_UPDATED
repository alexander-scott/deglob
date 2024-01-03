cc_library(
    name = "files_with_glob",
    deps = [":files_with_glob_file_1", ":files_with_glob_file_2"],
)

cc_library(
    name = "files_with_glob_file_1",
    hdrs = ["file_1.cpp"],
)

cc_library(
    name = "files_with_glob_file_2",
    hdrs = ["file_2.cpp"],
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":files_with_glob"],
)
