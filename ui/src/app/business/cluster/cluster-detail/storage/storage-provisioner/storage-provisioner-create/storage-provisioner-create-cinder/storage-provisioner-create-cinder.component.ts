import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageProvisionerService} from "../../storage-provisioner.service";
import {CreateStorageProvisionerRequest} from "../../storage-provisioner";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";
import {AlertLevels} from '../../../../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
    selector: 'app-storage-provisioner-create-cinder',
    templateUrl: './storage-provisioner-create-cinder.component.html',
    styleUrls: ['./storage-provisioner-create-cinder.component.css']
})
export class StorageProvisionerCreateCinderComponent implements OnInit {

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
        this.nfsForm.resetForm(this.item);
        this.item.name = 'cinder.csi.openstack.org';
        this.item.vars['enable_blockstorage'] = 'disable';
    }

    reset() {
        this.item = new CreateStorageProvisionerRequest();
    }

    changeBlockEnable() {
        if(this.item.vars['enable_blockstorage'] === 'disable') {
            this.item.vars['cinder_blockstorage_version'] = '';
            this.item.vars['node_volume_attach_limit'] = '';
        } else {
            this.item.vars['cinder_blockstorage_version'] = 'V3';
            this.item.vars['node_volume_attach_limit'] = '256';
        }
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
