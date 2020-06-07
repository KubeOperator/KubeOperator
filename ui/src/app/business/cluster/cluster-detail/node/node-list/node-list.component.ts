import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {KubernetesService} from "../../../kubernetes.service";
import {V1Namespace, V1Node} from "@kubernetes/client-node";
import {Cluster} from "../../../cluster";
import {ActivatedRoute} from "@angular/router";

@Component({
    selector: 'app-node-list',
    templateUrl: './node-list.component.html',
    styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {

    loading = true;
    selected = [];
    items: V1Node[] = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';
    page = 1;
    currentCluster: Cluster;
    @Output() openDetail = new EventEmitter<V1Node>();

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster.item;
            this.list();
        });
    }

    list() {
        this.loading = true;
        this.service.listNodes(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

    getInternalIp(n: V1Node) {
        let result = '';
        for (const addr of n.status.addresses) {
            if (addr.type === 'InternalIP') {
                result = addr.address;
            }
        }
        return result;
    }

    formatRAM(memory: string): string {
        let result = 0.0;
        if (memory.endsWith('Ki')) {
            const str = memory.substring(0, memory.indexOf('Ki'));
            result = parseFloat(str);
            result = result / (1024 * 1024);
        }
        return result.toFixed(2) + 'GB';
    }

    getNodeRoles(item: V1Node): string[] {
        const roles: string[] = [];
        for (const key in item.metadata.labels) {
            if (key) {
                switch (key) {
                    case 'node-role.kubernetes.io/master':
                        roles.push('master');
                        break;
                    case 'node-role.kubernetes.io/etcd':
                        roles.push('etcd');
                        break;
                }
            }
        }
        return roles;
    }

    isNodeReady(n: V1Node): string {
        let result = 'NotReady';
        for (const condition of n.status.conditions) {
            if (condition.type === 'Ready') {
                if (condition.status === 'True') {
                    result = 'Ready';
                }
            }
        }
        return result;
    }

    onDetail(item: V1Node) {
        this.openDetail.emit(item);
    }
}
