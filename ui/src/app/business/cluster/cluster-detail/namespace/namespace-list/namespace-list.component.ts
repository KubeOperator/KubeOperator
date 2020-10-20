import {Component, OnInit} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {V1Namespace} from '@kubernetes/client-node';
import {CommonAlertService} from '../../../../../layout/common-alert/common-alert.service';
import {TranslateService} from '@ngx-translate/core';
import {AlertLevels} from '../../../../../layout/app-alert/alert';
import {V1ObjectMeta} from '@kubernetes/client-node/dist/gen/model/v1ObjectMeta';

@Component({
    selector: 'app-namespace-list',
    templateUrl: './namespace-list.component.html',
    styleUrls: ['./namespace-list.component.css']
})
export class NamespaceListComponent implements OnInit {

    loading = true;
    selected = [];
    items: V1Namespace[] = [];
    page = 1;
    currentCluster: Cluster;
    opened = false;
    isSubmitGoing = false;
    namespace: string;

    constructor(private service: KubernetesService, private route: ActivatedRoute,
                private commonAlertService: CommonAlertService,
                private translateService: TranslateService) {
    }


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listNamespaces(this.currentCluster.name).subscribe(data => {
            this.loading = false;
            this.items = data.items;
        });
    }

    onCreate() {
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
        this.isSubmitGoing = false;
    }

    onSubmit() {
        const item = this.newV1NameSpace();
        this.isSubmitGoing = true;
        this.service.createNamespace(this.currentCluster.name, item).subscribe(res => {
            this.opened = false;
            this.isSubmitGoing = false;
            this.commonAlertService.showAlert(this.translateService.instant('APP_ADD_SUCCESS'), AlertLevels.SUCCESS);
            this.list();
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
