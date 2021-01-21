import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {StorageProvisionerListComponent} from './storage-provisioner-list/storage-provisioner-list.component';
import {StorageProvisionerCreateComponent} from './storage-provisioner-create/storage-provisioner-create.component';
import {StorageProvisionerCreateNfsComponent} from './storage-provisioner-create/storage-provisioner-create-nfs/storage-provisioner-create-nfs.component';
import {CreateStorageProvisionerRequest, StorageProvisioner} from './storage-provisioner';
import {StorageProvisionerDeleteComponent} from "./storage-provisioner-delete/storage-provisioner-delete.component";
import {StorageProvisionerCreateExternalCephComponent} from "./storage-provisioner-create/storage-provisioner-create-external-ceph/storage-provisioner-create-external-ceph.component";
import {StorageProvisionerCreateRookCephComponent} from "./storage-provisioner-create/storage-provisioner-create-rook-ceph/storage-provisioner-create-rook-ceph.component";
import {StorageProvisionerCreateVsphereComponent} from "./storage-provisioner-create/storage-provisioner-create-vsphere/storage-provisioner-create-vsphere.component";
import {StorageProvisionerCreateOceanStorComponent} from './storage-provisioner-create/storage-provisioner-create-ocean-stor/storage-provisioner-create-ocean-stor.component';

@Component({
    selector: 'app-storage-provisioner',
    templateUrl: './storage-provisioner.component.html',
    styleUrls: ['./storage-provisioner.component.css']
})
export class StorageProvisionerComponent implements OnInit {

    constructor(private route: ActivatedRoute) {
    }

    @ViewChild(StorageProvisionerListComponent, {static: true})
    list: StorageProvisionerListComponent;

    @ViewChild(StorageProvisionerCreateComponent, {static: true})
    create: StorageProvisionerCreateComponent;

    @ViewChild(StorageProvisionerDeleteComponent, {static: true})
    delete: StorageProvisionerDeleteComponent;

    @ViewChild(StorageProvisionerCreateNfsComponent, {static: true})
    nfs: StorageProvisionerCreateNfsComponent;
    @ViewChild(StorageProvisionerCreateExternalCephComponent, {static: true})
    externalCeph: StorageProvisionerCreateExternalCephComponent;
    @ViewChild(StorageProvisionerCreateRookCephComponent, {static: true})
    rookCeph: StorageProvisionerCreateRookCephComponent;
    @ViewChild(StorageProvisionerCreateVsphereComponent, {static: true})
    vsphere: StorageProvisionerCreateVsphereComponent;
    @ViewChild(StorageProvisionerCreateOceanStorComponent, {static: true})
    oceanStor: StorageProvisionerCreateOceanStorComponent;

    currentCluster: Cluster;

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
        });
    }

    openCreate() {
        this.create.open();
    }

    openDelete(items: StorageProvisioner[]) {
        this.delete.open(items);
    }


    openSelected(item: CreateStorageProvisionerRequest) {
        item.vars = {};
        item.name = '';
        switch (item.type) {
            case 'nfs':
                this.nfs.open(item);
                break;
            case 'external-ceph':
                this.externalCeph.open(item);
                break;
            case 'rook-ceph':
                this.rookCeph.open(item);
                break;
            case 'vsphere':
                this.vsphere.open(item);
                break;
            case 'oceanstor':
                this.oceanStor.open(item);
                break;
        }
    }

    refresh() {
        this.list.refresh();
    }

}
