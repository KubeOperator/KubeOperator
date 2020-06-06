import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { K8sPaginationComponent } from './k8s-pagination.component';

describe('K8sPaginationComponent', () => {
  let component: K8sPaginationComponent;
  let fixture: ComponentFixture<K8sPaginationComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ K8sPaginationComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(K8sPaginationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
