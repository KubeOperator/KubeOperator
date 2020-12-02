import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {MultiClusterRepository} from "../multi-cluster-repository";
import {MultiClusterRepositoryService} from "../multi-cluster-repository.service";
import {ModalAlertService} from "../../../shared/common-component/modal-alert/modal-alert.service";
import {AlertLevels} from "../../../layout/common-alert/alert";

@Component({
    selector: 'app-multi-cluster-repository-delete',
    templateUrl: './multi-cluster-repository-delete.component.html',
    styleUrls: ['./multi-cluster-repository-delete.component.css']
})
export class MultiClusterRepositoryDeleteComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    @Output() deleted = new EventEmitter();
    items: MultiClusterRepository[] = [];

    ngOnInit(): void {
    }

    open(items: MultiClusterRepository[]) {
        this.items = items;
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.multiClusterRepositoryService.batch('delete', this.items).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.deleted.emit();
        }, err => {
            this.modalAlertService.showAlert(err.error.msg, AlertLevels.ERROR);
            this.isSubmitGoing = false;
        });
    }
}
