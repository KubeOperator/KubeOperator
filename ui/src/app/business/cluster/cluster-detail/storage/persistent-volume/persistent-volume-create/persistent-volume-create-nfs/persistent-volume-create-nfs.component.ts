import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {V1PersistentVolume} from '@kubernetes/client-node';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';
import {V1PersistentVolumeSpec} from '@kubernetes/client-node/dist/gen/model/v1PersistentVolumeSpec';
import {V1NFSVolumeSource} from '@kubernetes/client-node/dist/gen/model/v1NFSVolumeSource';
import {NgForm} from '@angular/forms';
import {KubernetesService} from '../../../../../kubernetes.service';
import {Cluster} from '../../../../../cluster';

@Component({
    selector: 'app-persistent-volume-create-nfs',
    templateUrl: './persistent-volume-create-nfs.component.html',
    styleUrls: ['./persistent-volume-create-nfs.component.css']
})
export class PersistentVolumeCreateNfsComponent implements OnInit {

    constructor(private kubernetesService: KubernetesService) {
    }

    opened = false;
    item: V1PersistentVolume = this.newNfsPv();
    accessMode = 'ReadWriteOnce';
    isSubmitGoing = false;

    @Output() created = new EventEmitter();
    @Input() currentCluster: Cluster;
    @ViewChild('pvForm') pvForm: NgForm;


    ngOnInit(): void {
    }

    open() {
        this.reset();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    reset() {
        this.item = this.newNfsPv();
        this.pvForm.resetForm();
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.isSubmitGoing = true;
        this.item.spec.accessModes.push(this.accessMode);
        this.kubernetesService.createPersistentVolume(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.opened = false;
        });
    }

    newNfsPv(): V1PersistentVolume {
        return {
            apiVersion: 'v1',
            kind: 'PersistentVolume',
            metadata: {
                name: ''
            } as V1ObjectMeta,
            spec: {
                capacity: {},
                accessModes: [],
                nfs: {
                    path: '',
                    server: '',
                } as V1NFSVolumeSource,
            } as V1PersistentVolumeSpec,
        } as V1PersistentVolume;
    }
}

