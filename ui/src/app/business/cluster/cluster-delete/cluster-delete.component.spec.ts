import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterDeleteComponent } from './cluster-delete.component';

describe('ClusterDeleteComponent', () => {
  let component: ClusterDeleteComponent;
  let fixture: ComponentFixture<ClusterDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
