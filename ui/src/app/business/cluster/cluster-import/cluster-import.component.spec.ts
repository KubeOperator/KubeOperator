import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ClusterImportComponent } from './cluster-import.component';

describe('ClusterImportComponent', () => {
  let component: ClusterImportComponent;
  let fixture: ComponentFixture<ClusterImportComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ClusterImportComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ClusterImportComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
