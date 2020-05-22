import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterCreateComponent } from './cluster-create.component';

describe('ClusterCreateComponent', () => {
  let component: ClusterCreateComponent;
  let fixture: ComponentFixture<ClusterCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
