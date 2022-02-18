import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {CreateStorageClassRequest} from '../../storage';
import {NgForm} from '@angular/forms';
import {V1Secret, V1StorageClass} from '@kubernetes/client-node';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';
import {StorageProvisionerService} from '../../storage-provisioner/storage-provisioner.service';
import {StorageProvisioner} from '../../storage-provisioner/storage-provisioner';
import {Cluster} from '../../../../cluster';
import {KubernetesService} from '../../../../kubernetes.service';
import {ModalAlertService} from '../../../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../../layout/common-alert/alert';

@Component({
    selector: 'app-storage-class-create',
    templateUrl: './storage-class-create.component.html',
    styleUrls: ['./storage-class-create.component.css']
})
export class StorageClassCreateComponent implements OnInit {

    constructor(private provisionerService: StorageProvisionerService,
                private kubernetesService: KubernetesService,
                private modalAlertService: ModalAlertService) {
    }

    opened = false;
    isSecretsExit = false;
    isSubmitGoing = false;
    item: V1StorageClass = this.newV1StorageClass();
    provisioner: StorageProvisioner = new StorageProvisioner();
    @Input() currentCluster: Cluster;
    @Output() created = new EventEmitter();
    provisioners: StorageProvisioner[] = [];
    @ViewChild('itemForm') itemForm: NgForm;

    @Output() selected = new EventEmitter<CreateStorageClassRequest>();


    ngOnInit(): void {
    }

    reset() {
        this.isSubmitGoing = false;
        this.itemForm.resetForm();
        this.item = this.newV1StorageClass();
        this.provisioner = null;
        this.provisionerService.list(this.currentCluster.name).subscribe(data => {
            this.provisioners = []
            for (const provisioner of data) {
                if (provisioner.status === "Running") {
                    this.provisioners.push(provisioner)
                }
            }
            this.provisioners.push({
                name: 'kubernetes.io/no-provisioner',
                type: 'local-storage',
                vars: {},
                status: 'Running'
            } as StorageProvisioner);
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
            switch (this.provisioner.type) {
                case 'rook-ceph':
                    this.item.parameters['clusterID'] = 'rook-ceph';
                    this.item.parameters['pool'] = 'replicapool';
                    this.item.parameters['imageFormat'] = '2';
                    this.item.parameters['imageFeatures'] = 'layering';
                    this.item.parameters['csi.storage.k8s.io/provisioner-secret-name'] = 'rook-csi-rbd-provisioner';
                    this.item.parameters['csi.storage.k8s.io/provisioner-secret-namespace'] = 'rook-ceph';
                    this.item.parameters['csi.storage.k8s.io/controller-expand-secret-name'] = 'rook-csi-rbd-provisioner';
                    this.item.parameters['csi.storage.k8s.io/controller-expand-secret-namespace'] = 'rook-ceph';
                    this.item.parameters['csi.storage.k8s.io/node-stage-secret-name'] = 'rook-csi-rbd-node';
                    this.item.parameters['csi.storage.k8s.io/node-stage-secret-namespace'] = 'rook-ceph';
                    this.item.parameters['csi.storage.k8s.io/fstype'] = 'ext4';
                    break;
                case 'vsphere':
                    this.item.parameters['datastore'] = this.provisioner.vars['datastore'];
                    this.item.parameters['storagePolicyName'] = 'vSAN Default Storage Policy';
                    this.item.parameters['storagePolicyType'] = 'BuiltIn';
                    break;
                case 'glusterfs':
                    this.item.parameters['secretNamespace'] = 'kube-system';
                    this.item.parameters['restauthenabled'] = 'true';
                    this.item.parameters['gidMin'] = '40000';
                    this.item.parameters['gidMax'] = '50000';
                    this.item.parameters['volumetype'] = 'replicate:3';
                    break;
                case 'cinder':
                    this.item.allowVolumeExpansion = true;
                    break;
            }
        }
    }

    onSubmit() {
        if (this.isSubmitGoing) {
            return;
        }
        if (this.provisioner.type === 'glusterfs') {
            const mySecret = this.NewV1Secrets();
            this.kubernetesService.createSecret(this.currentCluster.name, this.item.parameters['secretNamespace'], mySecret).subscribe(data => {
                if (this.item.parameters['restuserkey']) {
                    delete this.item.parameters['restuserkey'];
                }
                this.addStorageClass()
            }, error => {
                this.modalAlertService.showAlert(error.message, AlertLevels.ERROR);
                return;
            });
        } else {
            this.addStorageClass()
        }
    }

    checkSecrets(){
        this.kubernetesService.getSecretByName(this.currentCluster.name, this.item.parameters['secretName'], this.item.parameters['secretNamespace']).subscribe(data => {
            this.isSecretsExit = true;
        }, error => {
            this.isSecretsExit = false;
            return;
        });
    }

    onCancel() {
        this.opened = false;
    }

    addStorageClass() {
        this.isSubmitGoing = true;

        if (this.item.parameters['storagePolicyType']) {
            delete this.item.parameters['storagePolicyType'];
        }
        this.kubernetesService.createStorageClass(this.currentCluster.name, this.item).subscribe(data => {
            this.isSubmitGoing = false;
            this.created.emit();
            this.opened = false;
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.message, AlertLevels.ERROR);
        });
    }

    NewV1Secrets(): V1Secret {
        return { 
            apiVersion: 'v1',
            kind: 'Secret',
            metadata: {
                name: this.item.parameters['secretName'],
                namespace: this.item.parameters['secretNamespace'],
            },
            stringData: {
                key: this.item.parameters['restuserkey'],
            },
            type: 'kubernetes.io/glusterfs'
        } as V1Secret;
    }

    newV1StorageClass(): V1StorageClass {
        return {
            apiVersion: 'storage.k8s.io/v1',
            kind: 'StorageClass',
            metadata: {
                name: ''
            } as V1ObjectMeta,
            provisioner: '',
            parameters: {},
        } as V1StorageClass;
    }
}
