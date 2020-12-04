import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {MultiClusterRepositoryService} from "../../../multi-cluster-repository.service";
import {ModalAlertService} from "../../../../../shared/common-component/modal-alert/modal-alert.service";
import {MultiClusterRepository} from "../../../multi-cluster-repository";
import {AlertLevels} from "../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-multi-cluster-relation-delete',
    templateUrl: './multi-cluster-relation-delete.component.html',
    styleUrls: ['./multi-cluster-relation-delete.component.css']
})
export class MultiClusterRelationDeleteComponent implements OnInit {

    constructor(private multiClusterRepositoryService: MultiClusterRepositoryService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    itemNames = [];
    @Input() currentRepository: MultiClusterRepository;
    @Output() deleted = new EventEmitter();
    ngOnInit(): void {
    }

    open(clusterNames: string[]) {
        this.itemNames = [];
        this.opened = true;
        this.itemNames = clusterNames;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.multiClusterRepositoryService.deleteRelations(this.currentRepository.name, this.itemNames).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.deleted.emit();
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }

}
