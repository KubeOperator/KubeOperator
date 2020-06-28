import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {CLusterImportRequest} from '../cluster';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';

@Component({
    selector: 'app-cluster-import',
    templateUrl: './cluster-import.component.html',
    styleUrls: ['./cluster-import.component.css']
})
export class ClusterImportComponent implements OnInit {

    constructor(private clusterService: ClusterService) {
    }

    opened = false;
    item = new CLusterImportRequest();
    @Output() imported = new EventEmitter();
    @ViewChild('importForm') importForm: NgForm;

    ngOnInit(): void {
    }

    open() {
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.clusterService.import(this.item).subscribe(() => {
            this.imported.emit();
            this.opened = false;
        });
    }

}
