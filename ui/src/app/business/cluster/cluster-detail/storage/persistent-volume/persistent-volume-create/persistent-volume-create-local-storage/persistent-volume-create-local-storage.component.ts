import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {KubernetesService} from "../../../../../kubernetes.service";
import {
    V1LocalVolumeSource,
    V1NodeSelector,
    V1NodeSelectorTerm,
    V1PersistentVolume,
    V1StorageClass
} from "@kubernetes/client-node";
import {Cluster} from "../../../../../cluster";
import {NgForm} from "@angular/forms";
import {V1NodeSelectorRequirement} from "@kubernetes/client-node/dist/gen/model/v1NodeSelectorRequirement";
import {V1ObjectMeta} from "@kubernetes/client-node/dist/gen/model/v1ObjectMeta";
import {V1VolumeNodeAffinity} from "@kubernetes/client-node/dist/gen/model/v1VolumeNodeAffinity";
import {V1PersistentVolumeSpec} from "@kubernetes/client-node/dist/gen/model/v1PersistentVolumeSpec";
import {ModalAlertService} from "../../../../../../../shared/common-component/modal-alert/modal-alert.service";
import {AlertLevels} from "../../../../../../../layout/common-alert/alert";

@Component({
    selector: 'app-persistent-volume-create-local-storage',
    templateUrl: './persistent-volume-create-local-storage.component.html',
    styleUrls: ['./persistent-volume-create-local-storage.component.css']
})
export class PersistentVolumeCreateLocalStorageComponent implements OnInit {

    constructor(private kubernetesService: KubernetesService, private alertService: ModalAlertService) {
    }

    opened = false;
    item: V1PersistentVolume = this.newLocalPv();
    accessMode = 'ReadWriteOnce';
    isSubmitGoing = false;
    selectorKey = '';
    selectorValue = '';
    selectorOperation = 'In';
    storageClazz: V1StorageClass[] = [];
    storageClassName = '';


    @Output()
    created = new EventEmitter();
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
        this.kubernetesService.listStorageClass(this.currentCluster.name, '', true).subscribe(data => {
            this.storageClazz = data.items.filter((sc) => {
                return sc.provisioner === 'kubernetes.io/no-provisioner';
            });
        });
        this.item = this.newLocalPv();
        this.pvForm.resetForm();
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        this.item.spec.accessModes.push(this.accessMode);
        if (this.selectorKey && this.selectorOperation) {
            this.item.spec.nodeAffinity.required.nodeSelectorTerms[0].matchExpressions[0] = {
                key: this.selectorKey,
                operator: this.selectorOperation,
            } as V1NodeSelectorRequirement;
            if (this.selectorValue) {
                this.item.spec.nodeAffinity.required.nodeSelectorTerms[0].matchExpressions[0].values = this.selectorValue.split(',');
            }
        } else {
            delete this.item.spec['nodeAffinity'];
        }
        if (this.storageClassName) {
            this.item.spec.storageClassName = this.storageClassName;
        }
        this.item.spec.capacity['storage'] += 'Gi';
        this.isSubmitGoing = true;
        this.kubernetesService.createPersistentVolume(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.opened = false;
        }, err => {
            this.isSubmitGoing = false;
            this.alertService.showAlert(err.error.message, AlertLevels.ERROR);
        });
    }

    newLocalPv(): V1PersistentVolume {
        return {
            apiVersion: 'v1',
            kind: 'PersistentVolume',
            metadata: {
                name: ''
            } as V1ObjectMeta,
            spec: {
                capacity: {},
                accessModes: [],
                storageClassName: '',
                local: {
                    path: '',
                } as V1LocalVolumeSource,
                nodeAffinity: {
                    required: {
                        nodeSelectorTerms: [
                            {
                                matchExpressions: [
                                    {
                                        key: '',
                                        operator: '',
                                        values: [],
                                    } as V1NodeSelectorRequirement,
                                ] as V1NodeSelectorRequirement[]
                            } as V1NodeSelectorTerm
                        ] as V1NodeSelectorTerm[],
                    } as V1NodeSelector,
                } as V1VolumeNodeAffinity,
            } as V1PersistentVolumeSpec,
        } as V1PersistentVolume;
    }

}
