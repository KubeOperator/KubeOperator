import {Component, OnInit} from '@angular/core';
import {Cluster} from "../../../../cluster";
import {NetworkingV1beta1Ingress, V1Namespace, V1Service} from "@kubernetes/client-node";
import {KubernetesService} from "../../../../kubernetes.service";
import {ActivatedRoute} from "@angular/router";
import {ignoreDiagnostics} from "@angular/compiler-cli/src/ngtsc/typecheck/src/diagnostics";

@Component({
    selector: 'app-ingress-list',
    templateUrl: './ingress-list.component.html',
    styleUrls: ['./ingress-list.component.css']
})
export class IngressListComponent implements OnInit {

    currentCluster: Cluster;
    items: NetworkingV1beta1Ingress[] = [];
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
        this.service.listIngress(this.currentCluster.name, this.continueToken, this.namespace).subscribe(data => {
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

    formatHost(ingress: NetworkingV1beta1Ingress): string[] {
        const results: string[] = [];
        for (const rule of ingress.spec.rules) {
            results.push(rule.host);
        }
        return results;
    }

    formatService(ingress: NetworkingV1beta1Ingress): string[] {
        const results: string[] = [];
        for (const rule of ingress.spec.rules) {
            for (const path of rule.http.paths) {
                results.push(path.backend.serviceName);
            }
        }
        return results;
    }


}
