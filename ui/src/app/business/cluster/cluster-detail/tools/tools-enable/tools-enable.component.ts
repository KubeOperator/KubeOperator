import {Component, EventEmitter, Input, OnInit, Output, ViewChild} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ClusterTool} from '../tools';
import {V1StorageClass} from '@kubernetes/client-node';
import {NgForm} from '@angular/forms';
import {Cluster} from '../../../cluster';
import {ToolsService} from '../tools.service';
import {EsIndexPattern, EsIndexPatternHelper} from './../../../../../constant/pattern';
import {TranslateService} from '@ngx-translate/core';
import {PasswordPattern} from '../../../../../constant/pattern';
import {AlertLevels} from '../../../../../layout/common-alert/alert';
import {ModalAlertService} from '../../../../../shared/common-component/modal-alert/modal-alert.service';

@Component({
    selector: 'app-tools-enable',
    templateUrl: './tools-enable.component.html',
    styleUrls: ['./tools-enable.component.css']
})
export class ToolsEnableComponent implements OnInit {

    constructor(private kubernetesService: KubernetesService,
                private toolsService: ToolsService,
                private modalAlertService: ModalAlertService,
                private translateService: TranslateService,
    ) {}
    esIndexPattern = EsIndexPattern;
    esIndexPatternHelper = EsIndexPatternHelper;
    opened = false;
    isSubmitGoing = false;
    nodeNum = 0;
    isCorrect = true;
    item: ClusterTool = new ClusterTool();
    storageClazz: V1StorageClass[] = [];
    namespaceList: string[] = [];
    nodeList: string[] = [];
    buttonLoading = false;
    passwordPattern = PasswordPattern;
    @ViewChild('itemForm') itemForm: NgForm;
    @Output() enabled = new EventEmitter();
    @Input() currentCluster: Cluster;

    private validationStateMap: any = {
        password: true,
        rePassword: true,
        rePwdCheck: true,
    };


    ngOnInit(): void {
        this.listStorageClass();
        this.listNode();
        this.listNamespaces();
    }
    onSubmit() {
        this.buttonLoading = true;
        this.checkIsCorrect();
        if (this.isCorrect) {
            this.toolsService.enable(this.currentCluster.name, this.item).subscribe(data => {
                this.opened = false;
                this.enabled.emit();
                this.buttonLoading = false;
            }, error => {
                this.buttonLoading = false;
            });
        } else {
            this.buttonLoading = false;
        }
    }
    checkIsCorrect() {
        if (this.item.name === 'logging') {
            if (this.item.vars['elasticsearch.replicas'] > this.nodeNum) {
                this.isCorrect = false;
                this.modalAlertService.showAlert(this.translateService.instant('APP_EFK_CREATE_REPLICAS'), AlertLevels.ERROR);
                return;
            }
        }
    }
    onCancel() {
        this.opened = false;
    }

    reset() {
        this.itemForm.resetForm();
    }

    open(item: ClusterTool) {
        this.reset();
        this.opened = true;
        this.setDefaultVars(item);
        this.item = item;
    }
    listStorageClass() {
        this.kubernetesService.listStorageClass(this.currentCluster.name, '', true).subscribe(data => {
            this.storageClazz = data.items;
        });
    }
    listNamespaces() {
        this.kubernetesService.listNamespaces(this.currentCluster.name).subscribe(data => {
            this.namespaceList = [];
            data.items.forEach(item => {
                this.namespaceList.push(item.metadata.name);
            });
        });
    }
    listNode() {
        this.kubernetesService.listNodes(this.currentCluster.name).subscribe(data => {
            this.nodeList = [];
            data.items.forEach(item => {
                this.nodeList.push(item.metadata.name);
            });
            this.nodeNum = this.nodeList.length;
        });
    }

    getValidationState(key: string): boolean {
        return this.validationStateMap[key];
    }

    handleValidation(key) {
        const cont = this.itemForm.controls[key];
        if (cont && cont.invalid && !cont.hasError) {
            this.validationStateMap[key] = false;
            return;
        }
        this.validationStateMap[key] = true;
        const r = /^(?=.*\d)(?=.*[a-zA-Z])[\da-zA-Z~!@#$%^&*]{6,30}$/g;
        r.lastIndex = 0;
        if (!r.test(this.itemForm.form.get('password').value)) {
            this.itemForm.controls['password'].setErrors({passwordError: false});
            this.validationStateMap['password'] = false;
            return;
        } else {
            this.validationStateMap['password'] = true;
        }

        if (this.itemForm.form.get('rePassword').value !== null && this.itemForm.form.get('password').value !== this.itemForm.form.get('rePassword').value) {
            this.itemForm.controls[key].setErrors({rePwdError: false});
            this.validationStateMap['rePwdCheck'] = false;
        } else {
            this.validationStateMap['rePwdCheck'] = true;
            this.itemForm.controls['password'].setErrors(null);
            this.itemForm.controls['rePassword'].setErrors(null);
        }
    }

    setDefaultVars(item: ClusterTool) {
        switch (item.name) {
            case 'prometheus':
                item.vars = {
                    'server.retention': 10,
                    'server.persistentVolume.enabled': false,
                    'server.persistentVolume.size': 10,
                    'server.persistentVolume.storageClass': '',
                };
                break;
            case 'chartmuseum':
                item.vars = {
                    'persistence.enabled': false,
                    'env.open.DISABLE_API': false,
                    'persistence.storageClass': '',
                    'service.type': 'NodePort',
                    'persistence.size': 10,
                };
                break;
            case 'registry':
                item.vars = {
                    'persistence.enabled': false,
                    'persistence.storageClass': '',
                    'service.type': 'NodePort',
                    'persistence.size': 10,
                };
                break;
            case 'logging':
                item.vars = {
                    'elasticsearch.esJavaOpts.item': 1,
                    'elasticsearch.replicas': 1,
                    'elasticsearch.persistence.enabled': false,
                    'elasticsearch.volumeClaimTemplate.resources.requests.storage': 10,
                    'elasticsearch.volumeClaimTemplate.storageClassName': '',
                };
                break;
            case 'loki':
                item.vars = {
                    'loki.persistence.enabled': false,
                    'loki.persistence.size': 8,
                    'loki.persistence.storageClassName': '',
                    'promtail.dockerPath': this.currentCluster.spec.dockerStorageDir,
                };
                break;
            case 'grafana':
                item.vars = {
                    'persistence.enabled': false,
                    'persistence.size': 10,
                    'persistence.storageClassName': '',
                    'adminPassword': '',
                    'adminPasswordRe': '',
                };
                break;
            case 'kubeapps':
                item.vars = {
                    'postgresql.persistence.enabled': false,
                    'postgresql.persistence.size': 10,
                    'global.storageClass': '',
                };
                break;
            case 'dashboard':
                item.vars = {};
                break;
        }
    }

}
