import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {CreateStorageClassRequest} from '../../storage';
import {NgForm} from '@angular/forms';
import {V1StorageClass} from '@kubernetes/client-node';
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
            this.provisioners = data;
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
                    break;
            }
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
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert(error.message, AlertLevels.ERROR);
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
            parameters: {},
        } as V1StorageClass;
    }


}
