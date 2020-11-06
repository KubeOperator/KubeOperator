import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';

@Component({
    selector: 'app-k8s-pagination',
    templateUrl: './k8s-pagination.component.html',
    styleUrls: ['./k8s-pagination.component.css']
})
export class K8sPaginationComponent implements OnInit {

    page = 1;
    @Output() pageChange = new EventEmitter();
    @Input() continueToken = '';
    @Output() continueTokenChange = new EventEmitter();
    @Input() previousToken = '';
    @Input() nextToken = '';

    constructor() {
    }

    ngOnInit(): void {
    }

    onNext() {
        this.page++;
        this.previousToken = this.continueToken;
        this.continueToken = this.nextToken;
        this.continueTokenChange.emit(this.continueToken);
        this.pageChange.emit();
    }

    onPrevious() {
        this.page--;
        this.continueToken = this.previousToken;
        this.continueTokenChange.emit(this.continueToken);
        this.pageChange.emit();
    }
}
