import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {CreateStorageClassRequest} from '../../storage';
import {NgForm} from '@angular/forms';
import {V1PersistentVolume, V1StorageClass} from '@kubernetes/client-node';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';
import {V1NFSVolumeSource} from '@kubernetes/client-node/dist/gen/model/v1NFSVolumeSource';
import {V1PersistentVolumeSpec} from '@kubernetes/client-node/dist/gen/model/v1PersistentVolumeSpec';
import {StorageProvisionerService} from '../../storage-provisioner/storage-provisioner.service';
import {StorageProvisioner} from '../../storage-provisioner/storage-provisioner';
import {Cluster} from '../../../../cluster';
import {KubernetesService} from '../../../../kubernetes.service';

@Component({
    selector: 'app-storage-class-create',
    templateUrl: './storage-class-create.component.html',
    styleUrls: ['./storage-class-create.component.css']
})
export class StorageClassCreateComponent implements OnInit {

    constructor(private provisionerService: StorageProvisionerService, private kubernetesService: KubernetesService) {
    }

    opened = false;
    isSubmitGoing = false;
    item: V1StorageClass;
    provisioner: StorageProvisioner = new StorageProvisioner();
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();
    provisioners: StorageProvisioner[] = [];

    @Output() selected = new EventEmitter<CreateStorageClassRequest>();


    ngOnInit(): void {

    }

    reset() {
        this.item = this.newV1StorageClass();
        this.provisionerService.list(this.currentCluster.name).subscribe(data => {
            this.provisioners = data;
        });
    }

    open() {
        this.opened = true;
        this.reset();
    }

    onProvisionerChange() {
        this.item.provisioner = '';
        if (this.provisioner) {
            this.item.provisioner = this.provisioner.name;
        }
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.kubernetesService.createStorageClass(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.opened = false;
        });
    }

    onCancel() {
        this.opened = false;
    }

    newV1StorageClass(): V1StorageClass {
        return {
            apiVersion: 'storage.k8s.io/v1',
            kind: 'StorageClass',
            metadata: {
                name: ''
            } as V1ObjectMeta,
            provisioner: '',
        } as V1StorageClass;
    }


}
