import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {ClusterService} from '../cluster.service';
import {Cluster} from '../cluster';

@Component({
    selector: 'app-cluster-delete',
    templateUrl: './cluster-delete.component.html',
    styleUrls: ['./cluster-delete.component.css']
})
export class ClusterDeleteComponent implements OnInit {

    opened = false;
    items: Cluster[] = [];
    @Output() deleted = new EventEmitter();

    constructor(private service: ClusterService) {
    }


    ngOnInit(): void {
    }

    open(items: Cluster[]) {
        this.items = items;
        console.log(items);
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.service.batch('delete', this.items).subscribe(data => {
            this.deleted.emit();
            this.opened = false;
        });
    }

}
