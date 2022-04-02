import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Cluster} from '../../../../cluster';
import {V1StorageClass} from '@kubernetes/client-node';
import {KubernetesService} from '../../../../kubernetes.service';
import {ActivatedRoute} from '@angular/router';

@Component({
    selector: 'app-storage-class-list',
    templateUrl: './storage-class-list.component.html',
    styleUrls: ['./storage-class-list.component.css']
})
export class StorageClassListComponent implements OnInit {

    currentCluster: Cluster;
    items: V1StorageClass[] = [];
    loading = true;
    selected = [];
    nextToken = '';
    previousToken = '';
    continueToken = '';
    @Output() createEvent = new EventEmitter();
    @Output() deleteEvent = new EventEmitter<V1StorageClass>();

    constructor(private service: KubernetesService, private route: ActivatedRoute) {
    }

    ngOnInit(): void {
        this.route.parent.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.list();
        });
    }

    onCreate() {
        this.createEvent.emit();
    }

    onDelete(item) {
        this.deleteEvent.emit(item);
    }

    list() {
        this.loading = true;
        let search = {
            kind: "storageclasslist",
            cluster: this.currentCluster.name,
            continue: this.continueToken,
            limit: 10,
            namespace: "",
            name: "",
        }
        this.service.listResource(search).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

    refresh() {
        this.list();
    }

}
