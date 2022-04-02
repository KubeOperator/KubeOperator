import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {AlertLevels} from '../../../../../layout/app-alert/alert';
import {V1Namespace} from '@kubernetes/client-node';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {Cluster} from '../../../cluster';

@Component({
    selector: 'app-namespace-create',
    templateUrl: './namespace-create.component.html',
    styleUrls: ['./namespace-create.component.css']
})
export class NamespaceCreateComponent implements OnInit {

    opened = false;
    isSubmitGoing = false;
    namespace: string;
    @Output() created = new EventEmitter();
    @Input() currentCluster: Cluster;

    constructor(private service: KubernetesService, private route: ActivatedRoute,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }

    ngOnInit(): void {
    }


    open() {
        this.opened = true;
        this.isSubmitGoing = false;
    }

    onCancel() {
        this.opened = false;
        this.isSubmitGoing = false;
    }

    onSubmit() {
        const item = this.newV1NameSpace();
        this.isSubmitGoing = true;
        let create = {
            cluster: this.currentCluster.name,
            kind: "namespace",
            namespace: "",
            info: item,
        }
        this.service.createResourceNs(create).subscribe(data => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.created.emit();
        }, error => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.namespace = '';
            this.commonAlertService.showAlert(error.error.message, AlertLevels.ERROR);
        });
    }

    newV1NameSpace(): V1Namespace {
        return {
            apiVersion: 'v1',
            kind: 'Namespace',
            metadata: {
                name: this.namespace,
            } as V1ObjectMeta,
        };
    }
}
