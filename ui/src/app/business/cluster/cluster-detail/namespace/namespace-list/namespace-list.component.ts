import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {KubernetesService} from '../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../../cluster';
import {V1Namespace} from '@kubernetes/client-node';

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
    @Output() deleteEvent = new EventEmitter<string>();
    @Output() createEvent = new EventEmitter<string>();
    defaultNamespaces: string[] = ['default', 'kube-node-lease', 'kube-operator', 'kube-public', 'kube-system'];

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
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

    onDelete(item: V1Namespace) {
        this.deleteEvent.emit(item.metadata.name);
    }

    onCreate() {
        this.createEvent.emit();
    }

    checkNamespace(name: string): boolean {
        return !(this.defaultNamespaces.indexOf(name) > -1);
    }
}
