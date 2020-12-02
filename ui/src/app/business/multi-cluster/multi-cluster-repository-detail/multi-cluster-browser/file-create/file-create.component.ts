import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {MultiClusterRepository, TreeNode} from "../../../multi-cluster-repository";
import {NgForm} from "@angular/forms";
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";
import {CommonAlertService} from "../../../../../layout/common-alert/common-alert.service";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-file-create',
    templateUrl: './file-create.component.html',
    styleUrls: ['./file-create.component.css']
})
export class FileCreateComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private commonAlertService: CommonAlertService) {
    }

    opened = false;
    item: TreeNode = new TreeNode();
    isSubmitGoing = false;
    @Output() created = new EventEmitter();
    @ViewChild('itemForm') itemForm: NgForm;
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
        this.item.path += `/${this.item.name}`;
        this.multiClusterRepositoryService.createOrDeleteTreeNode(this.currentRepository.name, this.item).subscribe(data => {
            this.opened = false;
            this.created.emit();
            this.isSubmitGoing = false;
        }, error => {
            this.opened = false;
            this.commonAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

    reset() {
        this.item = new TreeNode();
        this.itemForm.resetForm();
    }

    open(item: TreeNode) {
        this.reset();
        this.opened = true;
        this.item = item;
    }

}
