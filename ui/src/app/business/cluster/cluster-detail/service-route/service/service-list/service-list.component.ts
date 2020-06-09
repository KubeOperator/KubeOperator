import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1beta1CronJob, V1Namespace, V1Service} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-service-list',
    templateUrl: './service-list.component.html',
    styleUrls: ['./service-list.component.css']
})
export class ServiceListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1Service[] = [];
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
            this.currentCluster = data.cluster.item;
            this.listNamespace();
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listService(this.currentCluster.name, this.continueToken, this.namespace).subscribe(data => {
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
