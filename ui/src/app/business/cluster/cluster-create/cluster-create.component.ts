import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {ClusterCreateRequest} from '../cluster';
import {NgForm} from '@angular/forms';
import {ClusterService} from '../cluster.service';
import {ClrWizard} from '@clr/angular';

@Component({
    selector: 'app-cluster-create',
    templateUrl: './cluster-create.component.html',
    styleUrls: ['./cluster-create.component.css']
})
export class ClusterCreateComponent implements OnInit {

    opened = false;
    item: ClusterCreateRequest = new ClusterCreateRequest();
    @ViewChild('wizard', {static: true}) wizard: ClrWizard;
    @ViewChild('clusterForm') clusterForm: NgForm;
    @Output() created = new EventEmitter();

    constructor(private service: ClusterService) {
    }

    ngOnInit(): void {
    }

    open() {
        this.item = new ClusterCreateRequest();
        this.wizard.reset();
        this.clusterForm.resetForm();
        this.opened = true;
    }

    onCancel() {
        this.opened = false;
    }

    onSubmit() {
        this.service.create(this.item).subscribe(data => {
            this.opened = false;
            this.created.emit();
        });
    }
}
