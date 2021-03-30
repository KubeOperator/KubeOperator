import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageProvisionerService} from "../../storage-provisioner.service";
import {CreateStorageProvisionerRequest} from "../../storage-provisioner";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";
import {AlertLevels} from '../../../../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
    selector: 'app-storage-provisioner-create-glusterfs',
    templateUrl: './storage-provisioner-create-glusterfs.component.html',
    styleUrls: ['./storage-provisioner-create-glusterfs.component.css']
})
export class StorageProvisionerCreateGlusterfsComponent implements OnInit {

    constructor(private storageProvisionerService: StorageProvisionerService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    item: CreateStorageProvisionerRequest = new CreateStorageProvisionerRequest();
    @Output() created = new EventEmitter();
    @Input() currentCluster: Cluster;
    @ViewChild('nfsForm') nfsForm: NgForm;

    ngOnInit(): void {
    }

    open(item: CreateStorageProvisionerRequest) {
        this.opened = true;
        this.item = item;
        this.item.name = 'kubernetes.io/glusterfs';
        this.nfsForm.resetForm(this.item);
    }

    reset() {
        this.item = new CreateStorageProvisionerRequest();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.storageProvisionerService.create(this.currentCluster.name, this.item).subscribe(data => {
            this.opened = false;
            this.created.emit();
        }, error => {
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
