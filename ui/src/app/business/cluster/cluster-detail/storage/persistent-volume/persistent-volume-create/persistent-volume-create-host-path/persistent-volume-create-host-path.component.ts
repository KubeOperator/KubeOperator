import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {KubernetesService} from '../../../../../kubernetes.service';
import {V1HostPathVolumeSource, V1PersistentVolume} from '@kubernetes/client-node';
import {Cluster} from '../../../../../cluster';
import {NgForm} from '@angular/forms';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';
import {V1NFSVolumeSource} from '@kubernetes/client-node/dist/gen/model/v1NFSVolumeSource';
import {V1PersistentVolumeSpec} from '@kubernetes/client-node/dist/gen/model/v1PersistentVolumeSpec';

@Component({
    selector: 'app-persistent-volume-create-host-path',
    templateUrl: './persistent-volume-create-host-path.component.html',
    styleUrls: ['./persistent-volume-create-host-path.component.css']
})
export class PersistentVolumeCreateHostPathComponent implements OnInit {

    constructor(private kubernetesService: KubernetesService) {
    }

    opened = false;
    item: V1PersistentVolume = this.newHostPathPv();
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
        this.item = this.newHostPathPv();
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

    newHostPathPv(): V1PersistentVolume {
        return {
            apiVersion: 'v1',
            kind: 'PersistentVolume',
            metadata: {
                name: ''
            } as V1ObjectMeta,
            spec: {
                capacity: {},
                accessModes: [],
                hostPath: {
                    path: '',
                } as V1HostPathVolumeSource,
            } as V1PersistentVolumeSpec,
        } as V1PersistentVolume;
    }


}
