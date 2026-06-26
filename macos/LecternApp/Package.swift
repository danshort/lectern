// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "LecternApp",
    platforms: [.macOS(.v13)],
    dependencies: [
        .package(path: "../OpenSpecKit"),
        .package(url: "https://github.com/apple/swift-markdown.git", from: "0.4.0"),
    ],
    targets: [
        .executableTarget(
            name: "LecternApp",
            dependencies: [
                "OpenSpecKit",
                .product(name: "Markdown", package: "swift-markdown"),
            ]
        ),
    ]
)
