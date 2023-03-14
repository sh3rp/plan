load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/sh3rp/plan
gazelle(name = "gazelle")

go_library(
    name = "plan",
    srcs = [
        "const.go",
        "plan.go",
        "plandb.go",
        "web.go",
    ],
    importpath = "github.com/sh3rp/plan",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_boltdb_bolt//:go_default_library",
        "@com_github_oklog_ulid//:go_default_library",
    ],
)
