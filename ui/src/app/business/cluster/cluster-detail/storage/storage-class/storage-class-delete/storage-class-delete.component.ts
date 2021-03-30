import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {StorageProvisionerService} from '../../storage-provisioner/storage-provisioner.service';
import {KubernetesService} from '../../../../kubernetes.service';
import {ModalAlertService} from '../../../../../../shared/common-component/modal-alert/modal-alert.service';
import {Cluster} from '../../../../cluster';
import {AlertLevels} from '../../../../../../layout/common-alert/alert';
import {CommonAlertService} from '../../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';

@Component({
    selector: 'app-storage-class-delete',
    templateUrl: './storage-class-delete.component.html',
    styleUrls: ['./storage-class-delete.component.css']
})
export class StorageClassDeleteComponent implements OnInit {


    opened = false;
    submitGoing = false;
    deleteName = '';
    deleteClass: any;
    @Input() currentCluster: Cluster;
    @Output() deleted = new EventEmitter();

    constructor(private provisionerService: StorageProvisionerService,
                private kubernetesService: KubernetesService,
                private commonAlertService: CommonAlertService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }


    open(item) {
        this.deleteClass = item;
        this.deleteName = item.metadata.name;
        this.opened = true;
    }


    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.submitGoing = true;

        if (this.deleteClass.provisioner === 'kubernetes.io/glusterfs') {
            const namespace = this.deleteClass.parameters.secretNamespace;
            const secretName = this.deleteClass.parameters.secretName;
            this.kubernetesService.deleteSecret(this.currentCluster.name, namespace, secretName).subscribe(res => {
                this.delete();
            }, error => {
                this.submitGoing = false;
                this.modalAlertService.showAlert('', AlertLevels.ERROR);
            });
        } else {
            this.delete();
        }
    }

    delete() {
        this.kubernetesService.deleteStorageClass(this.currentCluster.name, this.deleteName).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.submitGoing = false;
            this.deleted.emit();
        }, error => {
            this.submitGoing = false;
            this.modalAlertService.showAlert('', AlertLevels.ERROR);
        });
    }
}
