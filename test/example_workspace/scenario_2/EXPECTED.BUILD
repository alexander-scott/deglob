cc_library(
    name = "files_starting_with_1",
    deps = [":files_starting_with_1_src_file_1_cpp"],
    hdrs = ["src/file_1.h"],
)

cc_library(
    name = "files_starting_with_1_src_file_1_cpp",
    srcs = ["src/file_1.cpp"],
    hdrs = ["src/file_1.h"],
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [
        ":files_starting_with_1",
        ":files_starting_with_2",
    ],
)

cc_library(
    name = "files_starting_with_2",
    deps = [":files_starting_with_2_src_file_2_cpp", ":files_starting_with_2_src_file_22_cpp"],
)

cc_library(
    name = "files_starting_with_2_src_file_2_cpp",
    hdrs = ["src/file_2.cpp"],
)

cc_library(
    name = "files_starting_with_2_src_file_22_cpp",
    hdrs = ["src/file_22.cpp"],
)
