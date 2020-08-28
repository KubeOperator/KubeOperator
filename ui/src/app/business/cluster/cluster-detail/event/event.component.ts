import {Component, OnInit} from '@angular/core';
import {KubernetesService} from '../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';

@Component({
    selector: 'app-event',
    templateUrl: './event.component.html',
    styleUrls: ['./event.component.css']
})
export class EventComponent implements OnInit {

    loading = false;
    currentCluster: Cluster;
    namespaces;
    events;
    currentNamespace: string;

    constructor(private kubernetesService: KubernetesService,
                private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.loading = true;
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.kubernetesService.listNamespaces(this.currentCluster.name).subscribe(res => {
                this.namespaces = res.items;
                if (this.namespaces.length > 0) {
                    const namespace = this.namespaces[0];
                    this.currentNamespace = namespace.metadata.name;
                    this.listEvents(this.currentNamespace);
                }
            });
        });
    }

    listEvents(namespace: string) {
        this.loading = true;
        this.kubernetesService.listEventsByNamespace(this.currentCluster.name, namespace).subscribe(res => {
            this.events = res.items;
            this.loading = false;
        });
    }
}
