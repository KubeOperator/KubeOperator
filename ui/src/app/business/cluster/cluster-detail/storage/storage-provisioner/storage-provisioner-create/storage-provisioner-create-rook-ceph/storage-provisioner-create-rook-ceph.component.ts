import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {StorageProvisionerService} from "../../storage-provisioner.service";
import {CreateStorageProvisionerRequest} from "../../storage-provisioner";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";

@Component({
    selector: 'app-storage-provisioner-create-rook-ceph',
    templateUrl: './storage-provisioner-create-rook-ceph.component.html',
    styleUrls: ['./storage-provisioner-create-rook-ceph.component.css']
})
export class StorageProvisionerCreateRookCephComponent implements OnInit {

    constructor(private storageProvisionerService: StorageProvisionerService) {
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
        console.log(this.item);
        this.storageProvisionerService.create(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.opened = false;
            this.created.emit();
        });
    }


}
