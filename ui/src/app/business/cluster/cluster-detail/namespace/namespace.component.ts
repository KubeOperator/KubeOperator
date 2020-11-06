import {Component, OnInit, ViewChild} from '@angular/core';
import {NamespaceDeleteComponent} from './namespace-delete/namespace-delete.component';
import {NamespaceListComponent} from './namespace-list/namespace-list.component';
import {NamespaceCreateComponent} from './namespace-create/namespace-create.component';
import {ActivatedRoute} from '@angular/router';
import {Cluster} from '../../cluster';

@Component({
    selector: 'app-namespace',
    templateUrl: './namespace.component.html',
    styleUrls: ['./namespace.component.css']
})
export class NamespaceComponent implements OnInit {

    @ViewChild(NamespaceDeleteComponent, {static: true})
    delete: NamespaceDeleteComponent;

    @ViewChild(NamespaceListComponent, {static: true})
    list: NamespaceListComponent;

    @ViewChild(NamespaceCreateComponent, {static: true})
    create: NamespaceCreateComponent;

    currentCluster: Cluster;

    constructor(private route: ActivatedRoute) {
    }


    ngOnInit(): void {
        this.route.parent.data.subscribe(data => {
            this.currentCluster = data.cluster;
            this.list.refresh();
        });
    }

    onDelete(item) {
        this.delete.open(item);
    }

    refresh() {
        this.list.refresh();
    }

    onCreate() {
        this.create.open();
    }
}
