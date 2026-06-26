import SwiftUI
import OpenSpecKit

struct ContentView: View {
    @EnvironmentObject var model: AppModel

    var body: some View {
        NavigationSplitView {
            Sidebar()
                .navigationSplitViewColumnWidth(min: 220, ideal: 280)
        } detail: {
            DetailView()
        }
        .toolbar {
            ToolbarItem(placement: .navigation) {
                Button { model.openPanel() } label: { Label("Open", systemImage: "folder") }
            }
            ToolbarItem {
                Button { model.reload() } label: { Label("Reload", systemImage: "arrow.clockwise") }
                    .disabled(model.project == nil)
            }
        }
    }
}

struct Sidebar: View {
    @EnvironmentObject var model: AppModel

    var body: some View {
        Group {
            if let project = model.project {
                List(selection: $model.selection) {
                    ForEach(project.changes, id: \.name) { change in
                        Section(change.name) {
                            artifactRows(change)
                        }
                    }
                }
            } else {
                emptyState
            }
        }
        .frame(minWidth: 220)
    }

    @ViewBuilder
    private func artifactRows(_ change: Change) -> some View {
        if change.proposal.present {
            row(change, .proposal, "Proposal", "doc.text")
        }
        if change.design.present {
            row(change, .design, "Design", "pencil.and.outline")
        }
        if !change.specFiles.isEmpty {
            SpecsGroup(change: change)
        }
        if change.tasks.present {
            row(change, .tasks, "Tasks", "checklist")
        }
    }

    private func row(_ change: Change, _ kind: ArtifactKind, _ title: String, _ icon: String) -> some View {
        Label(title, systemImage: icon)
            .tag(ArtifactRef(changeName: change.name, kind: kind))
    }

    private var emptyState: some View {
        VStack(spacing: 12) {
            Image(systemName: "books.vertical")
                .font(.system(size: 40))
                .foregroundStyle(.secondary)
            Text("No project open").font(.headline)
            Button("Open Project…") { model.openPanel() }
            if let err = model.loadError {
                Text(err)
                    .font(.callout)
                    .foregroundStyle(.red)
                    .multilineTextAlignment(.center)
                    .padding(.horizontal)
            }
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
        .padding()
    }
}

struct SpecsGroup: View {
    let change: Change
    @State private var expanded = true

    var body: some View {
        DisclosureGroup(isExpanded: $expanded) {
            ForEach(change.specFiles, id: \.name) { sf in
                Label(sf.name, systemImage: "doc.plaintext")
                    .tag(ArtifactRef(changeName: change.name, kind: .specFile(sf.name)))
            }
        } label: {
            Label("Specs", systemImage: "folder")
        }
    }
}

struct DetailView: View {
    @EnvironmentObject var model: AppModel

    var body: some View {
        if let ref = model.selection, let change = model.currentChange() {
            let artifact = model.artifact(for: ref, in: change)
            ScrollView {
                VStack(alignment: .leading, spacing: 14) {
                    if isSpec(ref.kind), !artifact.readError {
                        ValidationBanner(change: change)
                    }
                    content(artifact)
                }
                .frame(maxWidth: .infinity, alignment: .leading)
                .padding(24)
            }
            .navigationTitle(change.name)
        } else {
            ContentUnavailableLikeView()
        }
    }

    @ViewBuilder
    private func content(_ artifact: Artifact) -> some View {
        if artifact.readError {
            Label(artifact.content, systemImage: "exclamationmark.triangle.fill")
                .foregroundStyle(.orange)
                .font(.callout)
        } else if !artifact.present {
            Text("This artifact is not present.").foregroundStyle(.secondary)
        } else {
            MarkdownView(artifact.content)
        }
    }

    private func isSpec(_ kind: ArtifactKind) -> Bool {
        if case .specFile = kind { return true }
        return false
    }
}

struct ValidationBanner: View {
    let change: Change

    var body: some View {
        let issues = validateChange(change)
        if !issues.isEmpty {
            VStack(alignment: .leading, spacing: 4) {
                Label("Validation issues", systemImage: "exclamationmark.triangle.fill")
                    .font(.callout.bold())
                    .foregroundStyle(.orange)
                ForEach(issues, id: \.self) { issue in
                    Text("• \(issue)").font(.callout)
                }
            }
            .padding(12)
            .frame(maxWidth: .infinity, alignment: .leading)
            .background(Color.orange.opacity(0.12))
            .clipShape(RoundedRectangle(cornerRadius: 8))
        }
    }
}

// Minimal stand-in (ContentUnavailableView is macOS 14+; keep deployment at 13).
struct ContentUnavailableLikeView: View {
    var body: some View {
        VStack(spacing: 8) {
            Image(systemName: "doc.text.magnifyingglass")
                .font(.system(size: 36))
                .foregroundStyle(.secondary)
            Text("Select an artifact").foregroundStyle(.secondary)
        }
        .frame(maxWidth: .infinity, maxHeight: .infinity)
    }
}
