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
    @Output() deleteEvent = new EventEmitter<string>();

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

    onDelete(name) {
        this.deleteEvent.emit(name);
    }

    list() {
        this.loading = true;
        this.service.listStorageClass(this.currentCluster.name, this.continueToken).subscribe(data => {
            this.loading = false;
            this.items = data.items;
            this.nextToken = data.metadata[this.service.continueTokenKey] ? data.metadata[this.service.continueTokenKey] : '';
        });
    }

    refresh() {
        this.list();
    }

}
