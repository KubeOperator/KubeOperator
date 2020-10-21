import {Component, OnInit, ViewChild} from '@angular/core';
import {NamespaceDeleteComponent} from './namespace-delete/namespace-delete.component';
import {NamespaceListComponent} from './namespace-list/namespace-list.component';
import {NamespaceCreateComponent} from './namespace-create/namespace-create.component';

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

    constructor() {
    }

    ngOnInit(): void {
        this.list.list();
    }

    onDelete(item) {
        this.delete.open(item);
    }

    refresh() {
        this.list.list();
    }

    onCreate() {
        this.create.open();
    }
}
