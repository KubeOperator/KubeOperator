import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1ConfigMap, V1Namespace, V1Secret} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-secret-list',
    templateUrl: './secret-list.component.html',
    styleUrls: ['./secret-list.component.css']
})
export class SecretListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1Secret[] = [];
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
        this.service.listSecret(this.currentCluster.name, this.continueToken, this.namespace).subscribe(data => {
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
