import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {MultiClusterRepository, TreeNode} from "../../../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {NgForm} from "@angular/forms";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-file-delete',
    templateUrl: './file-delete.component.html',
    styleUrls: ['./file-delete.component.css']
})
export class FileDeleteComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private commonAlertService: CommonAlertService) {
    }

    opened = false;
    item: TreeNode = new TreeNode();
    isSubmitGoing = false;
    @Output() deleted = new EventEmitter();
    @Input() currentRepository: MultiClusterRepository;

    ngOnInit(): void {
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.item.delete = true;
        this.multiClusterRepositoryService.createOrDeleteTreeNode(this.currentRepository.name, this.item).subscribe(data => {
            this.opened = false;
            this.deleted.emit();
            this.isSubmitGoing = false;
        }, error => {
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }



    open(item: TreeNode) {
        this.opened = true;
        this.item = item;
    }
}
