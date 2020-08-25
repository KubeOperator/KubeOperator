import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ManifestListComponent } from './manifest-list.component';

describe('ManifestListComponent', () => {
  let component: ManifestListComponent;
  let fixture: ComponentFixture<ManifestListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ManifestListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ManifestListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
