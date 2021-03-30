import {Component, OnInit, ViewChild} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {StorageProvisionerListComponent} from './storage-provisioner-list/storage-provisioner-list.component';
import {StorageProvisionerSyncComponent} from './storage-provisioner-sync/storage-provisioner-sync.component';
import {StorageProvisionerCreateComponent} from './storage-provisioner-create/storage-provisioner-create.component';
import {StorageProvisionerCreateNfsComponent} from './storage-provisioner-create/storage-provisioner-create-nfs/storage-provisioner-create-nfs.component';
import {CreateStorageProvisionerRequest, StorageProvisioner} from './storage-provisioner';
import {StorageProvisionerDeleteComponent} from "./storage-provisioner-delete/storage-provisioner-delete.component";
import {StorageProvisionerCreateExternalCephComponent} from "./storage-provisioner-create/storage-provisioner-create-external-ceph/storage-provisioner-create-external-ceph.component";
import {StorageProvisionerCreateRookCephComponent} from "./storage-provisioner-create/storage-provisioner-create-rook-ceph/storage-provisioner-create-rook-ceph.component";
import {StorageProvisionerCreateVsphereComponent} from "./storage-provisioner-create/storage-provisioner-create-vsphere/storage-provisioner-create-vsphere.component";
import {StorageProvisionerCreateOceanStorComponent} from './storage-provisioner-create/storage-provisioner-create-ocean-stor/storage-provisioner-create-ocean-stor.component';
import {StorageProvisionerCreateCinderComponent} from "./storage-provisioner-create/storage-provisioner-create-cinder/storage-provisioner-create-cinder.component";
import {StorageProvisionerCreateGlusterfsComponent} from "./storage-provisioner-create/storage-provisioner-create-glusterfs/storage-provisioner-create-glusterfs.component";

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

    @ViewChild(StorageProvisionerSyncComponent, {static: true})
    sync: StorageProvisionerSyncComponent;

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
    @ViewChild(StorageProvisionerCreateCinderComponent, {static: true})
    cinder: StorageProvisionerCreateCinderComponent;
    @ViewChild(StorageProvisionerCreateGlusterfsComponent, {static: true})
    glusterfs: StorageProvisionerCreateGlusterfsComponent;

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
            case 'cinder':
                this.cinder.open(item);
                break;
            case 'glusterfs':
                this.glusterfs.open(item);
                break;
        }
    }

    openSync (items: StorageProvisioner[]) {
        this.sync.open(items);
    }

    refresh() {
        this.list.refresh();
    }

}
