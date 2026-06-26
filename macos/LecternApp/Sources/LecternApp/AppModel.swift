import SwiftUI
import AppKit
import OpenSpecKit

// Identifies a single artifact within a change, used as the sidebar selection.
enum ArtifactKind: Hashable {
    case proposal
    case design
    case tasks
    case specFile(String)
}

struct ArtifactRef: Hashable {
    let changeName: String
    let kind: ArtifactKind
}

@MainActor
final class AppModel: ObservableObject {
    @Published var project: Project?
    @Published var rootPath: String?
    @Published var loadError: String?
    @Published var selection: ArtifactRef?

    private let bookmarkKey = "projectBookmark"
    private var accessedURL: URL?

    init() {
        restoreBookmark()
    }

    // MARK: - Opening / restoring the project (security-scoped bookmark)

    func openPanel() {
        let panel = NSOpenPanel()
        panel.canChooseDirectories = true
        panel.canChooseFiles = false
        panel.allowsMultipleSelection = false
        panel.message = "Choose a project folder containing an openspec/ directory"
        panel.prompt = "Open"
        guard panel.runModal() == .OK, let url = panel.url else { return }
        persistBookmark(for: url)
        load(url)
    }

    func reload() {
        if let path = rootPath {
            load(URL(fileURLWithPath: path))
        }
    }

    private func persistBookmark(for url: URL) {
        if let data = try? url.bookmarkData(options: .withSecurityScope,
                                            includingResourceValuesForKeys: nil, relativeTo: nil) {
            UserDefaults.standard.set(data, forKey: bookmarkKey)
        }
    }

    private func restoreBookmark() {
        guard let data = UserDefaults.standard.data(forKey: bookmarkKey) else { return }
        var stale = false
        guard let url = try? URL(resolvingBookmarkData: data, options: .withSecurityScope,
                                 relativeTo: nil, bookmarkDataIsStale: &stale) else { return }
        load(url)
        if stale { persistBookmark(for: url) }
    }

    private func load(_ url: URL) {
        accessedURL?.stopAccessingSecurityScopedResource()
        _ = url.startAccessingSecurityScopedResource()
        accessedURL = url
        rootPath = url.path
        do {
            let loaded = try Loader().loadFrom(url.path)
            project = loaded
            loadError = nil
            selection = defaultSelection(for: loaded)
        } catch {
            project = nil
            selection = nil
            loadError = describe(error)
        }
    }

    private func defaultSelection(for project: Project) -> ArtifactRef? {
        guard let first = project.changes.first else { return nil }
        if first.proposal.present { return ArtifactRef(changeName: first.name, kind: .proposal) }
        return ArtifactRef(changeName: first.name, kind: .design)
    }

    private func describe(_ error: Error) -> String {
        switch error {
        case LoaderError.noOpenspecDir(let root):
            return "No openspec/ directory found in \(root)"
        default:
            return "\(error)"
        }
    }

    // MARK: - Resolving the current selection

    func currentChange() -> Change? {
        guard let ref = selection else { return nil }
        return project?.changes.first { $0.name == ref.changeName }
    }

    func artifact(for ref: ArtifactRef, in change: Change) -> Artifact {
        switch ref.kind {
        case .proposal: return change.proposal
        case .design: return change.design
        case .tasks: return change.tasks
        case .specFile(let name):
            if let sf = change.specFiles.first(where: { $0.name == name }) {
                return Artifact(content: sf.content, present: true, readError: sf.readError)
            }
            return Artifact()
        }
    }
}
