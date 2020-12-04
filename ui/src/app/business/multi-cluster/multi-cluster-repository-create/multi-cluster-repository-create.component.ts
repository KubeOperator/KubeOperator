import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {MultiClusterRepositoryCreateRequest} from "../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../multi-cluster-repository.service";
import {ModalAlertService} from "../../../shared/common-component/modal-alert/modal-alert.service";
import {AlertLevels} from "../../../layout/common-alert/alert";
import {NgForm} from "@angular/forms";

@Component({
    selector: 'app-multi-cluster-repository-create',
    templateUrl: './multi-cluster-repository-create.component.html',
    styleUrls: ['./multi-cluster-repository-create.component.css']
})
export class MultiClusterRepositoryCreateComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    item: MultiClusterRepositoryCreateRequest = new MultiClusterRepositoryCreateRequest();
    @ViewChild('itemForm') itemForm: NgForm;
    @Output() created = new EventEmitter();

    ngOnInit(): void {
    }

    open() {
        this.reset();
        this.opened = true;
    }

    reset() {
        this.isSubmitGoing = false;
        this.item = new MultiClusterRepositoryCreateRequest();
        this.itemForm.resetForm();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.multiClusterRepositoryService.create(this.item).subscribe(data => {
            this.created.emit();
            this.isSubmitGoing = false;
            this.opened = false;
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
