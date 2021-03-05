import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageProvisionerService} from "../../storage-provisioner.service";
import {CreateStorageProvisionerRequest} from "../../storage-provisioner";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";
import {ModalAlertService} from '../../../../../../../shared/common-component/modal-alert/modal-alert.service';
import {AlertLevels} from '../../../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-storage-provisioner-create-vsphere',
    templateUrl: './storage-provisioner-create-vsphere.component.html',
    styleUrls: ['./storage-provisioner-create-vsphere.component.css']
})
export class StorageProvisionerCreateVsphereComponent implements OnInit {

    constructor(private storageProvisionerService: StorageProvisionerService,private modalAlertService: ModalAlertService) {
    }

    opened = false;
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
        this.item.name = 'csi.vsphere.vmware.com';
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
