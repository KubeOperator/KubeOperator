import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {Cluster} from '../../../../../cluster';
import {CreateStorageProvisionerRequest} from '../../storage-provisioner';
import {NgForm} from '@angular/forms';
import {StorageProvisionerService} from '../../storage-provisioner.service';
import {AlertLevels} from '../../../../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
    selector: 'app-storage-provisioner-create-ocean-stor',
    templateUrl: './storage-provisioner-create-ocean-stor.component.html',
    styleUrls: ['./storage-provisioner-create-ocean-stor.component.css']
})
export class StorageProvisionerCreateOceanStorComponent implements OnInit {

    item: CreateStorageProvisionerRequest = new CreateStorageProvisionerRequest();
    opened = false;
    isSubmitGoing = false;
    @ViewChild('storForm') storForm: NgForm;
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();

    constructor(private storageProvisionerService: StorageProvisionerService, private modalAlertService: ModalAlertService) {
    }

    ngOnInit(): void {
    }

    open(item: CreateStorageProvisionerRequest) {
        this.item = item;
        this.item.name = 'csi.huawei.com';
        this.opened = true;
        this.storForm.resetForm(this.item);
        this.item.vars['oceanstor_controller_type'] = 'single';
        this.item.vars['oceanstor_is_multipath'] = 'no';
    }

    onCancel() {
        this.opened = false;
        this.isSubmitGoing = false;
    }

    onSubmit() {
        this.isSubmitGoing = true;
        this.storageProvisionerService.create(this.currentCluster.name, this.item).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.created.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.error.msg, AlertLevels.ERROR);
        });
    }
}
