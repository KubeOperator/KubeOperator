import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {StorageProvisionerService} from '../../storage-provisioner/storage-provisioner.service';
import {KubernetesService} from '../../../../kubernetes.service';
import {CommonAlertService} from '../../../../../../layout/common-alert/common-alert.service';
import {ModalAlertService} from '../../../../../../shared/common-component/modal-alert/modal-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../../layout/common-alert/alert';
import {Cluster} from '../../../../cluster';

@Component({
    selector: 'app-persistent-volume-delete',
    templateUrl: './persistent-volume-delete.component.html',
    styleUrls: ['./persistent-volume-delete.component.css']
})
export class PersistentVolumeDeleteComponent implements OnInit {


    opened = false;
    isSubmitGoing = false;
    deleteName = '';
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


    open(name) {
        this.deleteName = name;
        this.opened = true;
    }


    onCancel() {
        this.opened = false;
    }


    onSubmit() {
        this.isSubmitGoing = true;
        this.kubernetesService.deletePersistentVolume(this.currentCluster.name, this.deleteName).subscribe(res => {
            this.commonAlertService.showAlert(this.translateService.instant('APP_DELETE_SUCCESS'), AlertLevels.SUCCESS);
            this.opened = false;
            this.isSubmitGoing = false;
            this.deleted.emit();
        }, error => {
            this.isSubmitGoing = false;
            this.modalAlertService.showAlert('', AlertLevels.ERROR);
        });
    }
}
