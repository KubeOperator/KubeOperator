import {Component, OnInit, ViewChild} from '@angular/core';
import {MultiClusterRepositoryService} from "../../multi-cluster-repository.service";
import {FileContent, MultiClusterRepository, TreeNode} from "../../multi-cluster-repository";
import {ActivatedRoute} from "@angular/router";
import {FileCreateComponent} from "./file-create/file-create.component";
import {FileDeleteComponent} from "./file-delete/file-delete.component";
import {CommonAlertService} from "../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../layout/common-alert/alert";

@Component({
    selector: 'app-multi-cluster-browser',
    templateUrl: './multi-cluster-browser.component.html',
    styleUrls: ['./multi-cluster-browser.component.css']
})


export class MultiClusterBrowserComponent implements OnInit {
    cmOptions: any = {
        lineNumbers: true,
        styleActiveLine: true,
        lineWrapping: true,
        // mode: {name: 'yaml'},
        theme: 'neo',
    };
    currentRepository: MultiClusterRepository;
    editFiles: TreeNode[] = [];
    tree: TreeNode;
    @ViewChild(FileCreateComponent, {static: true})
    fileCreateView: FileCreateComponent;
    @ViewChild(FileDeleteComponent, {static: true})
    fileDeleteView: FileDeleteComponent;
    remoteLoading = false;
    isSubmitGoing = false;

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService,
                private route: ActivatedRoute, private commonAlertService: CommonAlertService) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(d => {
            this.currentRepository = d.repo;
            this.refresh();
        });
    }

    getChildren = (folder: TreeNode) => folder.children;

    onOpenFile(file: TreeNode) {
        if (file.dir) {
            return;
        }
        for (const node of this.editFiles) {
            if (file.path === node.path) {
                return;
            }
        }
        if (file.changed) {
            this.editFiles.push(file);
            for (const node of this.editFiles) {
                node.active = node.name === file.name;
            }
            return;
        }
        this.multiClusterRepositoryService.readFile(this.currentRepository.name, file.path).subscribe(data => {
            file.content = data.content;
            file.originContent = data.content;
            this.editFiles.push(file);
            for (const node of this.editFiles) {
                node.active = node.name === file.name;
            }
        });
    }

    closeFile(file: TreeNode) {
        this.editFiles.splice(this.editFiles.indexOf(file), 1);
        for (const f of this.editFiles) {
            if (f.active) {
                return;
            }
        }
        if (this.editFiles.length > 0) {
            this.editFiles[this.editFiles.length - 1].active = true;
        }
    }


    onContentChange(file: TreeNode) {
        file.changed = file.content !== file.originContent;
    }

    onResetContent(file: TreeNode) {
        file.content = file.originContent;
    }


    openCreate(node: TreeNode, type: string) {
        const t = new TreeNode();
        if (type === 'dir') {
            t.dir = true;
        }
        t.path = node.path;
        this.fileCreateView.open(t);
    }

    openDelete(node: TreeNode) {
        this.fileDeleteView.open(node);
    }

    saveContent(node: TreeNode) {
        const content = new FileContent();
        content.path = node.path;
        content.content = node.content;
        this.multiClusterRepositoryService.saveFile(this.currentRepository.name, content).subscribe(data => {
            node.changed = false;
        });
    }

    pullRemoteRepository() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.remoteLoading = true;
        this.multiClusterRepositoryService.pullRemoteRepository(this.currentRepository.name).subscribe(data => {
            this.isSubmitGoing = false;
            this.remoteLoading = false;
            this.refresh();
        }, error => {
            this.isSubmitGoing = false;
            this.remoteLoading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    pushRemoteRepository() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.remoteLoading = true;
        this.multiClusterRepositoryService.pushRemoteRepository(this.currentRepository.name).subscribe(data => {
            this.isSubmitGoing = false;
            this.remoteLoading = false;
            this.refresh();
        }, error => {
            this.isSubmitGoing = false;
            this.remoteLoading = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    refresh() {
        this.editFiles = [];
        this.multiClusterRepositoryService.getTree(this.currentRepository.name).subscribe(data => {
            const root = new TreeNode();
            root.path = '/';
            root.name = '/';
            root.dir = true;
            root.children.push(data);
            this.tree = root;
        }, error => {
        });
    }

}
