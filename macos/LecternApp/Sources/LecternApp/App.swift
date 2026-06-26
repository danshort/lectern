import SwiftUI
import AppKit

// LecternApp is the native macOS reader for OpenSpec artifacts. It reuses the
// OpenSpecKit domain layer (the Swift port of Go internal/openspec) and renders
// with swift-markdown. Per the App Sandbox decision (Option C), it is
// non-sandboxed but accesses the project through a security-scoped bookmark from
// day one, so a future sandbox flip changes nothing here.

final class AppDelegate: NSObject, NSApplicationDelegate {
    func applicationDidFinishLaunching(_ notification: Notification) {
        // A SwiftPM executable launches as an accessory by default; promote it
        // to a regular app so it gets a Dock icon, menu bar, and key focus.
        NSApp.setActivationPolicy(.regular)
        NSApp.activate(ignoringOtherApps: true)
    }

    func applicationShouldTerminateAfterLastWindowClosed(_ sender: NSApplication) -> Bool { true }
}

@main
struct LecternApp: App {
    @NSApplicationDelegateAdaptor(AppDelegate.self) private var delegate
    @StateObject private var model = AppModel()

    var body: some Scene {
        WindowGroup("lectern") {
            ContentView()
                .environmentObject(model)
                .frame(minWidth: 820, minHeight: 520)
        }
        .commands {
            CommandGroup(replacing: .newItem) {
                Button("Open Project…") { model.openPanel() }
                    .keyboardShortcut("o", modifiers: .command)
            }
        }
    }
}
