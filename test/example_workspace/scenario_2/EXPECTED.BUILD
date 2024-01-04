cc_library(
    name = "files_starting_with_1",
    deps = [":files_starting_with_1_file_1"],
    hdrs = ["file_1.h"],
)

cc_library(
    name = "files_starting_with_1_file_1",
    srcs = ["file_1.cpp"],
    hdrs = ["file_1.h"],
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
    deps = [":files_starting_with_2_file_2", ":files_starting_with_2_file_22"],
)

cc_library(
    name = "files_starting_with_2_file_2",
    hdrs = ["file_2.cpp"],
)

cc_library(
    name = "files_starting_with_2_file_22",
    hdrs = ["file_22.cpp"],
)
