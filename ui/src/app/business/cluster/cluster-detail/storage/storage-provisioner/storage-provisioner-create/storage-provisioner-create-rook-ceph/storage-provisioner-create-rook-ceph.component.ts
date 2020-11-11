import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageProvisionerService} from "../../storage-provisioner.service";
import {CreateStorageProvisionerRequest} from "../../storage-provisioner";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";
import {ModalAlertService} from '../../../../../../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-storage-provisioner-create-rook-ceph',
    templateUrl: './storage-provisioner-create-rook-ceph.component.html',
    styleUrls: ['./storage-provisioner-create-rook-ceph.component.css']
})
export class StorageProvisionerCreateRookCephComponent implements OnInit {

    constructor(private storageProvisionerService: StorageProvisionerService, private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSubmitGoing = false;
    item: CreateStorageProvisionerRequest = new CreateStorageProvisionerRequest();
    @Output() created = new EventEmitter();
    @Input() currentCluster: Cluster;
    @ViewChild('nfsForm') nfsForm: NgForm;

    ngOnInit(): void {
    }

    open(item: CreateStorageProvisionerRequest) {
        this.reset();
        this.opened = true;
        this.item = item;
        this.item.name = 'rook-ceph.rbd.csi.ceph.com';
        this.nfsForm.resetForm(this.item);
    }

    reset() {
        this.item = new CreateStorageProvisionerRequest();
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.storageProvisionerService.create(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.created.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }


}
