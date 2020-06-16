import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1DaemonSet, V1Deployment, V1Namespace} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-daemon-set-list',
    templateUrl: './daemon-set-list.component.html',
    styleUrls: ['./daemon-set-list.component.css']
})
export class DaemonSetListComponent implements OnInit {
    currentCluster: Cluster;
    items: V1DaemonSet[] = [];
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
        this.service.listDaemonSet(this.currentCluster.name, this.continueToken, this.namespace).subscribe(data => {
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
