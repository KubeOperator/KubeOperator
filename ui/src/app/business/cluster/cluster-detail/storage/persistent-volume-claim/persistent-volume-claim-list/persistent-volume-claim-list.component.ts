import {Component, OnInit} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1Namespace, V1PersistentVolume} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-persistent-volume-claim-list',
    templateUrl: './persistent-volume-claim-list.component.html',
    styleUrls: ['./persistent-volume-claim-list.component.css']
})
export class PersistentVolumeClaimListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1PersistentVolume[] = [];
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
        });
    }

    list() {
        this.loading = true;
        let search = {
            kind: "pvlist",
            cluster: this.currentCluster.name,
            continue: this.continueToken,
            limit: 10,
            namespace: this.namespace,
            name: "",
        }
        this.service.listResource(search).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

    listNamespace() {
        this.loading = true;
        let search = {
            kind: "namespacelist",
            cluster: this.currentCluster.name,
            continue: "",
            limit: 0,
            namespace: "",
            name: "",
        }
        this.service.listResource(search).subscribe(data => {
            this.namespaces = data.items;
            if (this.namespace === '') {
                this.namespace = this.items[0].metadata.name;
            }
        });
        this.list();
    }

}
