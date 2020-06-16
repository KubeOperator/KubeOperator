import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1beta1CronJob, V1Namespace} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-corn-job-list',
    templateUrl: './corn-job-list.component.html',
    styleUrls: ['./corn-job-list.component.css']
})
export class CornJobListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1beta1CronJob[] = [];
    namespaces: V1Namespace[] = [];
    namespace = '';
    loading = true;
    selected = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster =data.cluster;
            this.listNamespace();
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listCornJob(this.currentCluster.name, this.continueToken, this.namespace).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

    listNamespace() {
        this.service.listNamespaces(this.currentCluster.name).subscribe(data => {
            this.namespaces = data.items;
        });
    }
}
